<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="end">
      <el-button v-show="hasManagerPermission()" type="primary" icon="el-icon-plus" @click="handleAdd">添加</el-button>
    </el-row>
    <el-table
      border
      stripe
      highlight-current-row
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="160" />
      <el-table-column prop="name" label="组名" />
      <el-table-column prop="insertTime" label="创建时间" width="160" />
      <el-table-column prop="updateTime" label="更新时间" width="160" />
      <el-table-column prop="operation" label="操作" width="150">
        <template slot-scope="scope">
          <el-button size="small" type="primary" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button :disabled="!hasGroupManagerPermission()" size="small" type="danger" @click="handleRemove(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-row type="flex" justify="end" style="margin-top: 10px;">
      <el-pagination
        hide-on-single-page
        :total="pagination.total"
        :page-size="pagination.rows"
        background
        layout="prev, pager, next"
        @current-change="handlePageChange"
      />
    </el-row>
    <el-dialog title="分组设置" :visible.sync="dialogVisible">
      <el-form ref="form" :rules="formRules" :model="formData" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button :disabled="formProps.disabled" type="primary" @click="submit">确 定</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import { getList, getTotal, add, edit, remove } from '@/api/group'

export default {
  data() {
    return {
      dialogVisible: false,
      tableData: [],
      pagination: {
        page: 1,
        rows: 16,
        total: 0
      },
      tempFormData: {},
      formProps: {
        disabled: false
      },
      formData: {
        id: 0
      },
      formRules: {
        name: [
          { required: true, message: '请输入分组名称', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.storeFormData()
    this.getList()
    this.getTotal()
  },
  methods: {
    getList() {
      getList(this.pagination).then((response) => {
        this.tableData = response.data.list
      })
    },

    getTotal() {
      getTotal().then((response) => {
        this.pagination.total = response.data.total
      })
    },

    // 分页事件
    handlePageChange(val) {
      this.pagination.page = val
      this.getList()
    },

    handleAdd() {
      this.restoreFormData()
      this.dialogVisible = true
    },

    handleEdit(data) {
      this.formData = Object.assign({}, data)
      this.dialogVisible = true
    },

    handleRemove(data) {
      this.$confirm('此操作将删除该分组, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        remove(data.id).then((response) => {
          this.$message.success('删除成功')
          this.getList()
          this.getTotal()
        })
      }).catch(() => {
        this.$message.info('已取消删除')
      })
    },

    submit() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          if (this.formData.id === 0) {
            this.add()
          } else {
            this.edit()
          }
        } else {
          return false
        }
      })
    },

    add() {
      this.formProps.disabled = true
      add(this.formData).then((response) => {
        this.getList()
        this.getTotal()
        this.$message.success('添加成功')
      }).finally(() => {
        this.formProps.disabled = this.dialogVisible = false
      })
    },

    edit() {
      this.formProps.disabled = true
      edit(this.formData).then((response) => {
        this.getList()
        this.$message.success('修改成功',)
      }).finally(() => {
        this.formProps.disabled = this.dialogVisible = false
      })
    },

    storeFormData() {
      this.tempFormData = JSON.parse(JSON.stringify(this.formData))
    },

    restoreFormData() {
      this.formData = JSON.parse(JSON.stringify(this.tempFormData))
    }
  }
}
</script>
